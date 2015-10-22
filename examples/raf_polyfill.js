/* jslint brower: true */

(function() {
  'use strict';

  var duration = 1000.0 / 60;

  window.requestAnimationFrame = (
    window.requestAnimationFrame ||
    window.msRequestAnimationFrame ||
    window.mozRequestAnimationFrame ||
    window.webkitRequestAnimationFrame ||
    window.oRequestAnimationFrame ||
    function(callback) {
      return setTimeout(function() {
        callback(new Date().getTime());
      }, duration);
    }
  );

  window.cancelAnimationFrame = (
    window.cancelAnimationFrame ||
    window.msCancelAnimationFrame ||
    window.mozCancelAnimationFrame ||
    window.webkitCancelAnimationFrame ||
    window.oCancelAnimationFrame ||
    function(id) {
      clearTimeout(id);
    }
  );
}());
