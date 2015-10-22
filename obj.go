package goggles

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

type ObjGroup struct {
	Name         string
	MaterialName string
	IndexBuffer  *js.Object
	NumIndices   int
	faces        [][]string
	indices      []uint16
}

type Obj struct {
	Name         string
	Groups       []*ObjGroup
	VertexBuffer *js.Object
	materials    map[string]*js.Object
	tupleIndex   uint16
	tupleIndices map[string]uint16
	vertices     []float32
}

func (o *Obj) GetMaterial(name string) *js.Object {
	return o.materials[name]
}

func (o *Obj) SetMaterial(name string, texture *js.Object) {
	if o.materials == nil {
		o.materials = map[string]*js.Object{}
	}
	o.materials[name] = texture
}

func (o *Obj) Read(reader io.Reader, gl *webgl.Context) error {
	var group *ObjGroup
	var materialName string
	var positions []float32
	var normals []float32
	var texcoords []float32

	var float2 [2]float32 = [2]float32{0, 0}
	var float3 [3]float32 = [3]float32{0, 0, 0}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			// Blank line or comment, ignore it
			continue
		}

		// Split line into fields on whitespace
		fields := strings.Fields(line)
		switch strings.ToLower(fields[0]) {
		// Vertex position.
		case "v":
			if err := parseFloat3(fields[1:4], &float3); err != nil {
				return err
			}
			positions = append(positions, float3[0], float3[1], float3[2])

		// Vertex normal.
		case "vn":
			if err := parseFloat3(fields[1:4], &float3); err != nil {
				return err
			}
			normals = append(normals, float3[0], float3[1], float3[2])

		// Vertex texture coordinates.
		case "vt":
			if err := parseFloat2(fields[1:3], &float2); err != nil {
				return err
			}
			texcoords = append(texcoords, float2[0], 1.0-float2[1])

		// Face indices, specified in sets of "position/uv/normal".
		case "f":
			faces := fields[1:len(fields)]
			if group == nil {
				group = &ObjGroup{MaterialName: materialName}
				o.Groups = append(o.Groups, group)
			}
			group.faces = append(group.faces, faces)

		// New group, with a name.
		case "g":
			group = &ObjGroup{Name: fields[1], MaterialName: materialName}
			o.Groups = append(o.Groups, group)

		// Object name. The obj will only have one object statement.
		case "o":
			o.Name = fields[1]

		// Material library. I'm not handling this for now. Instead, call
		// SetMaterial() for each of the named materials.
		// case "mtllib":

		// Specifies the material for the current group (and any future groups
		// that don't have their own usemtl statement).
		case "usemtl":
			materialName = fields[1]
			if group != nil {
				group.MaterialName = materialName
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Normalize the vertices. OBJ does two things that cause problems for
	// modern renderers: it allows faces to be polygons, instead of only
	// triangles; and it allows each face vertex to have a different index
	// for each stream (position, normal, uv).
	//
	// This code creates triangle fans out of any faces that have more than
	// three vertexes, and merges distinct groupings of pos/normal/uv into
	// a single vertex stream.
	var faceIndices [3]uint16 = [3]uint16{0, 0, 0}
	o.tupleIndex = 0
	for _, g := range o.Groups {
		for _, f := range g.faces {
			for i := 1; i < len(f)-1; i++ {
				var err error
				if faceIndices[0], err = o.mergeTuple(f[i], positions, normals, texcoords); err != nil {
					return err
				}
				if faceIndices[1], err = o.mergeTuple(f[0], positions, normals, texcoords); err != nil {
					return err
				}
				if faceIndices[2], err = o.mergeTuple(f[i+1], positions, normals, texcoords); err != nil {
					return err
				}
				g.indices = append(g.indices, faceIndices[0], faceIndices[1], faceIndices[2])
			}
		}
		g.IndexBuffer = gl.CreateBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g.IndexBuffer)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, g.indices, gl.STATIC_DRAW)
		g.NumIndices = len(g.indices)
	}

	o.VertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, o.VertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, o.vertices, gl.STATIC_DRAW)

	return nil
}

func (o *Obj) mergeTuple(tuple string, p []float32, n []float32, t []float32) (uint16, error) {
	var err error

	if index, ok := o.tupleIndices[tuple]; ok {
		return index, nil
	}

	indexStrings := strings.Split(tuple, "/")

	// Position index
	var pi uint64
	if pi, err = strconv.ParseUint(indexStrings[0], 10, 16); err != nil {
		return 0, err
	}
	pi--
	o.vertices = append(o.vertices, p[pi*3+0], p[pi*3+1], p[pi*3+2])

	// Normal index
	if len(indexStrings) > 2 && indexStrings[2] != "" {
		var ni uint64
		if ni, err = strconv.ParseUint(indexStrings[2], 10, 16); err != nil {
			return 0, err
		}
		ni--
		o.vertices = append(o.vertices, n[ni*3+0], n[ni*3+1], n[ni*3+2])
	} else {
		// Face doesn't have a normal
		o.vertices = append(o.vertices, 0, 0, 0)
	}

	// Texcoord index
	if len(indexStrings) > 1 && indexStrings[1] != "" {
		var ti uint64
		if ti, err = strconv.ParseUint(indexStrings[1], 10, 16); err != nil {
			return 0, err
		}
		ti--
		o.vertices = append(o.vertices, t[ti*2+0], t[ti*2+1])
	} else {
		// Face doesn't have a texcoord
		o.vertices = append(o.vertices, 0, 0)
	}

	// Cache the merged tuple index in case it's used again
	tupleIndex := o.tupleIndex
	o.tupleIndex++
	if o.tupleIndices == nil {
		o.tupleIndices = map[string]uint16{}
	}
	o.tupleIndices[tuple] = tupleIndex
	return tupleIndex, nil
}

func parseFloat2(strings []string, floats *[2]float32) error {
	var f float64
	var err error
	for i := 0; i < 2; i++ {
		f, err = strconv.ParseFloat(strings[i], 32)
		if err != nil {
			return err
		}
		floats[i] = float32(f)
	}
	return nil
}

func parseFloat3(strings []string, floats *[3]float32) error {
	var f float64
	var err error
	for i := 0; i < 3; i++ {
		f, err = strconv.ParseFloat(strings[i], 32)
		if err != nil {
			return err
		}
		floats[i] = float32(f)
	}
	return nil
}
