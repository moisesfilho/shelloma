#!/usr/bin/env python3
import base64, os, re, sys

svg_path = "scripts/org.shelloma.Shelloma.svg"
with open(svg_path) as f:
    svg = f.read()

m = re.search(r'href="data:image/png;base64,([^"]+)"', svg)
png = base64.b64decode(m.group(1))

for arg in sys.argv[1:]:
    sz, path = arg.split("=", 1)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "wb") as f:
        f.write(png)
