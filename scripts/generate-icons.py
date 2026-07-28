#!/usr/bin/env python3
import base64, io, os, re, sys
from PIL import Image

SVG_PATH = "scripts/org.shelloma.Shelloma.svg"

def extract_png_from_svg(svg_path):
    with open(svg_path) as f:
        content = f.read()
    m = re.search(r'href="data:image/png;base64,([^"]+)"', content)
    if not m:
        print(f"✖ Nenhum PNG embutido em {svg_path}", file=sys.stderr)
        sys.exit(1)
    return Image.open(io.BytesIO(base64.b64decode(m.group(1)))).convert("RGBA")

def embed_png_in_svg(img, svg_path):
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    b64 = base64.b64encode(buf.getvalue()).decode("ascii")
    svg = f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <image href="data:image/png;base64,{b64}" width="512" height="512"/>
</svg>'''
    with open(svg_path, "w") as f:
        f.write(svg)
    print(f"✔ SVG regenerado: {svg_path} ({len(svg) // 1024} KB)")

if __name__ == "__main__":
    if not os.path.exists(SVG_PATH):
        print(f"✖ {SVG_PATH} não encontrado", file=sys.stderr)
        sys.exit(1)

    img = extract_png_from_svg(SVG_PATH)
    img_512 = img.resize((512, 512), Image.LANCZOS)

    sizes = []
    for a in sys.argv[1:]:
        if a.isdigit():
            sizes.append((int(a), f"{a}.png"))
        elif "=" in a:
            s, path = a.split("=", 1)
            sizes.append((int(s), path))

    for s, path in sizes:
        out = img_512.resize((s, s), Image.LANCZOS)
        out.save(path)
        print(f"✔ PNG {s}x{s} -> {path}")

    if "--regenerate" in sys.argv:
        embed_png_in_svg(img_512, SVG_PATH)
