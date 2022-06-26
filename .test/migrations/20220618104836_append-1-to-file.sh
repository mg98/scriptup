### migrate up ###
echo 1 >> .test/foo.txt
### migrate down ###
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  sed -i '$ d' .test/foo.txt
else # macos
  sed -i '' -e '$ d' .test/foo.txt
fi