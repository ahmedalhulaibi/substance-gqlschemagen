cd model
gengostruct -jsonsrc='../substance-objects.json'
cd ..
gengqlschema
gqlgen