# svg2scad

Command-line tool that converts SVG paths into bezier curves for use in [OpenSCAD](https://openscad.org).
It uses the [BOSL2 library](https://github.com/BelfrySCAD/BOSL2) to represent the curves, resulting in an OpenSCAD module
that has the nice features of BOSL2 like [attachability](https://github.com/BelfrySCAD/BOSL2/wiki/attachments.scad).