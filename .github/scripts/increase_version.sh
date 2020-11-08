set -eux
git fetch --tags || echo "could not fetch tags"
latest_version=$(git describe --tags --match "[0-9]*.[0-9]*.[0-9]*" --abbrev=0)

current_branch=$(git branch --show-current)

if [[ $current_branch != feature/* ]] && [[ $current_branch != patch/* ]];
then
    echo "no feature or patch branch found"
    exit 0
fi

if [[ $current_branch == feature/* ]];
then
    echo "increasing minor version"
    new_version=$(python3 .github/scripts/version.py ${latest_version} --minor)
fi

if [[ $current_branch == patch/* ]];
then
    echo "increasing patch version"
    new_version=$(python3 .github/scripts/version.py ${latest_version} --patch)
fi

git config --global user.email "mdreem@fastmail.fm"
git config --global user.name "Github Action"

git tag -a ${new_version} -m ${new_version}
git push origin ${new_version}
