### Descriptor language

Each descriptor item is define in between '{{' and '}}'.

1. @model <model_name>      //  Possible Values: cql | mongo | redis
This is the first directive which defines the Message as a model.

2. @tab ( \[ (<partition_key>, ...), \] <clustering_key>, ...)
3. @view ( \[ (<partition_key>, ...), \] <clustering_key>, ...)
4. @counter < 