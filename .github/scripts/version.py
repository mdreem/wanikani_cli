import argparse

parser = argparse.ArgumentParser(description="Increase the version.")
parser.add_argument(
    "version_string", metavar="<VER>", type=str, help="the version string"
)
parser.add_argument(
    "--minor", dest="minor", action="store_true", help="increase the minor version"
)
parser.add_argument(
    "--patch", dest="patch", action="store_true", help="increase the patch version"
)

args = parser.parse_args()

major, minor, patch = args.version_string.split(".")

if args.patch:
    patch = int(patch) + 1

if args.minor:
    minor = int(minor) + 1
    patch = 0

print(f"{major}.{minor}.{patch}")
