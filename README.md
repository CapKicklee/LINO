# LINO : Large Input, Narrow Output

 LINO is a simple ETL (Extract Transform Load) tools to manage tests datas.
 The `lino` command line tool
extract test data from a relational database to create a smallest
production-like database.

## Usage

`lino` command line work in relative project's directory, like `git` or `docker`


## Create a new LINO project

```
$ mkdir myproject
$ cd myproject
```

## Add DataConnector

A DataConnector is a database connection shortcut.

```bash
$ lino dataconnector add source postgresql://postgres:sakila@localhost:5432/postgres?sslmode=disable
successfully added dataconnector {source postgresql://postgres:sakila@localhost:5432/postgres?sslmode=disable}
```

The content of `dataconnector.yml` generated is

```yaml
version: v1
dataconnectors:
  - name: source
    url: postgresql://postgres:sakila@localhost:5432/postgres?sslmode=disable
```


## Create relationships

LINO create a consistent sample database. To perform extraction that respect foreign keys constraints LINO have to extract relationships between tables.

Use the `relationships` sub-command or its short name `rs` to extract relationship from foreign key constraints.


```
$ lino relation extract source
lino finds 40 relations from constraints
```

The content of `relationships.yml` generated is

```yaml
version: v1
relations:
  - name: film_original_language_id_fkey
    parent:
        name: public.film
        keys:
          - original_language_id
    child:
        name: public.language
        keys:
          - language_id
  - name: film_language_id_fkey
.
.
.
```

At least user can edit the `relationships.yml` manually to add personal relationship.

## Ingress descriptor

Ingress descriptor object describe how `lino` has to go through the relationships to extract data test.

### Create Ingress descriptor

To create ingress descriptor use the `ingress-descriptor` sub-command (or its shortcut `id`) with the start table of the extraction.

```bash
$ lino id create public.customer
successfully created ingress descriptor
```

`lino` store the new ingress descriptor in `ingress-descriptor.yml` file :

```yaml
version: v1
IngressDescriptor:
    startTable: public.customer
    relations:
      - name: film_original_language_id_fkey
        parent:
            name: public.film
            lookup: false
        child:
            name: public.language
```

### Display plan

The `display-plan` utilities explain the `lino`'s plan to extract data from database.

```bash
$ lino id display-plan
step 1 - extract rows from public.customer
step 2 - extract rows from public.store following →customer_store_id_fkey relationship for rows extracted at step 1, then follow →store_manager_staff_id_fkey →staff_store_id_fkey relationships (loop until data exhaustion)
step 3 - extract rows from public.address following →staff_address_id_fkey relationship for rows extracted at step 2
step 4 - extract rows from public.city following →address_city_id_fkey relationship for rows extracted at step 3
.
.
.
```

### Show graph

The `show-graph` create a graph of tables as node and relation as edge.

```bash
$ lino id customer show-graph
```

`lino` open your browser to visualize graph generated.

![Test Image 1](doc/img/lino-graph-export.svg)

## Table

The `lino table` action extract informations about tables.

```
$ lino table extract source
lino finds 15 table(s)
```

`lino` store the table description in `table.yml` file :

```yaml
version: v1
tables:
  - name: public.actor
    keys:
      - actor_id
  - name: public.address
    keys:
      - address_id
  - name: public.category
    keys:
```

## Extract

The `extract` sub-command create à **json** object for each line (jsonline format http://jsonlines.org/) of the first table.

```
$ lino extract source
{"active":1,"activebool":true,"address_id":5,"create_date":"2006-02-14T00:00:00Z","customer_address_id_fkey":{"address":"1913 Hanoi Way","address2":"","address_city_id_fkey":{"city":"Sasebo","city_country_id_fkey":{"country":"Japan","country_id":50,"last_update":"2006-02-15T09:44:00Z"},"city_id":463,"country_id":50,"last_update":"2006-02-15T09:45:25Z"},"address_id":5,"city_id":463,"district":"Nagasaki","last_update":"2006-02-15T09:45:30Z","phone":"28303384290","postal_code":"35200"},"customer_id":1,"customer_store_id_fkey":{"address_id":1,"last_update":"2006-02-15T09:57:12Z","manager_staff_id":1,"store_address_id_fkey":{"address":"47 MySakila Drive","address2":null,"address_city_id_fkey":{"city":"Lethbridge","city_country_id_fkey":{"country":"Canada","country_id":20,"last_update":"2006-02-15T09:44:00Z"},"city_id":300,"country_id":20,"last_update":"2006-02-15T09:45:25Z"},"address_id":1,"city_id":300,"district":"Alberta","last_update":"2006-02-15T09:45:30Z","phone":"","postal_code":""},"store_id":1,"store_manager_staff_id_fkey":{"active":true,"address_id":3,"email":"Mike.Hillyer@sakilastaff.com","first_name":"Mike","last_name":"Hillyer","last_update":"2006-05-16T16:13:11.79328Z","password":"8cb2237d0679ca88db6464eac60da96345513964","picture":"iVBORw0KWgo=","staff_address_id_fkey":{"address":"23 Workhaven Lane","address2":null,"address_city_id_fkey":{"city":"Lethbridge","city_country_id_fkey":{"country":"Canada","country_id":20,"last_update":"2006-02-15T09:44:00Z"},"city_id":300,"country_id":20,"last_update":"2006-02-15T09:45:25Z"},"address_id":3,"city_id":300,"district":"Alberta","last_update":"2006-02-15T09:45:30Z","phone":"14033335568","postal_code":""},"staff_id":1,"store_id":1,"username":"Mike"}},"email":"MARY.SMITH@sakilacustomer.org","first_name":"MARY","last_name":"SMITH","last_update":"2006-02-15T09:57:20Z","store_id":1}
```

## Load

TODO

### Interaction with other tools

**LINO** respect the UNIX philosophy and use standards input an output to share data with others tools.

## MongoDB storage

Data set could be store in mongoDB easily with the `mongoimport` tool:

```
$ lino extract source --limit 100 | mongoimport --db myproject --collection customer
```

and reload later to a database :

```bash
$ mongoexport --db myproject --collection customer | lino load customer --jdbc jdbc:oracle:thin:scott/tiger@target:1721:xe
```

## jq

`jq` tool can be piped with the **LINO** output to prettify it.

```
$ lino extract source | jq
```

Extract sub field from the JSON stream

```
$ lino extract source --limit 3 | jq ".email"
"MARY.SMITH@sakilacustomer.org"
"PATRICIA.JOHNSON@sakilacustomer.org"
"LINDA.WILLIAMS@sakilacustomer.org"
```

Project subfield to produce other *JSON* objects

```
$ lino extract source --limit 3 | jq '{ "manager": .customer_store_id_fkey.store_manager_staff_id_fkey.first_name , "customer_email" :  .email }'

{
  "manager": "Mike",
  "customer_email": "MARY.SMITH@sakilacustomer.org"
}
{
  "manager": "Mike",
  "customer_email": "PATRICIA.JOHNSON@sakilacustomer.org"
}
{
  "manager": "Mike",
  "customer_email": "LINDA.WILLIAMS@sakilacustomer.org"
}
```

## Installation

Download the last binary release in your path

### Contribute

TODO