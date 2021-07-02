### Descriptor language

Each descriptor item is defined in between '{{' and '}}'.
You don't need to define them directly, instead you could use Protobuf options defined
in `options.proto` in the root of the project.

1. @model <model_name>        
> Possible Values for <model_name>: cql | local
2. @tab ( \[ (<partition_key>, ...), \] <clustering_key>, ...)
3. @view ( \[ (<partition_key>, ...), \] <clustering_key>, ...)



### Commands
`goda graph ./... | dot -Tsvg -o ./docs/dec_graph.svg`