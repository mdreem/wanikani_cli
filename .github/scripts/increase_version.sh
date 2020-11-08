set -eux
git fetch --tags || echo "could not fetch tags"
latest_version=$(git describe --tags --match "[0-9]*.[0-9]*.[0-9]*" --abbrev=0)

merged_message=$(git log -1 | grep "Merge pull request .* from mdreem/.*" -o) \
    || (echo "not a properly formatted merge commit" && exit 0)

if [[ "${merged_message}" =~ Merge\ pull\ request\ .*\ from\ mdreem/(.*) ]];
then
    merged_branch=${BASH_REMATCH[1]}
    echo "merged branch: ${merged_branch}"
else
    echo "not a properly formatted merge commit"
    exit 0
fi

if [[ $merged_branch != feature/* ]] && [[ $merged_branch != patch/* ]];
then
    echo "no feature or patch branch found"
    exit 0
fi

if [[ $merged_branch == feature/* ]];
then
    echo "increasing minor version"
    new_version=$(python3 .github/scripts/version.py ${latest_version} --minor)
fi

if [[ $merged_branch == patch/* ]];
then
    echo "increasing patch version"
    new_version=$(python3 .github/scripts/version.py ${latest_version} --patch)
fi

git config --global user.email "mdreem@fastmail.fm"
git config --global user.name "Github Action"

git tag -a ${new_version} -m ${new_version}
git push origin ${new_version}
