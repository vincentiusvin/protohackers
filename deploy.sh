OUT="./build/$(basename $1 .go)"

echo Building $1 at $OUT
go build -o "$OUT" $1
echo Built!
echo Copying...
scp "$OUT" root@vincentiusvin.com:/protohackers/
echo Copied!
