# `/teacher`

- api Version: `v1`

## POST - `/new`

### Request

NOTE: Data must be submited as _Form Data_!

* `name`: The Name of the Teacher; _string_
* `email`: Mail Adress of the Teacher; _string_

### Response

* `200 OK`: Teacher (and referenced User) was sucessfully created
* `4xx`: There was an Error while creating the Teacher/ User

### Example

```bash
~] curl -X POST http://localhost:3000/v1/teacher/new \
	-F 'name=Max Mustermann' \
	-F 'email=max.mustermann@rbs-ulm.de'
```
