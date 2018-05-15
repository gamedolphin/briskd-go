name=briskd
rm $name
go build -o $name ../src

./$name -color $@

