## Encrypt secrets in CUE using SOPS


Configure an age key and sops configuration file:

```bash
AGE_KEY_NAME="cuesops.txt"
AGE_CONF_DIR="$HOME/.config/sops/age"
AGE_KEY_PATH=${AGE_CONF_DIR}/${AGE_KEY_NAME}

mkdir -p ${AGE_CONF_DIR}

rm -f ${AGE_KEY_PATH}
age-keygen -o ${AGE_KEY_PATH}

AGE=$(awk '/public/{print $4}' ${AGE_KEY_PATH})

cat > .sops.yaml <<EOF
creation_rules:
- age: ${AGE}
EOF
```

We use annotations on the CUE values that we wish to encrypt:

```cue
data: {
	api_key: "this is an api key" @secret(sops) // <-- will be encrypted
}

SECRET: "this is a secret value" @secret(sops) // <-- will be encrypted

another: { // <-- will not be encrypted
	this:  100
	value: "is not encrypted"
}
```

Save the above example to a file named `secrets.cue`.

Run the following to encrypt the values in `secrets.cue`:

```
go run main.go encrypt secrets.cue
```


Now our cue file looks like this:

```cue
data: {
	api_key: "ENC[AES256_GCM,data:0SeH+BIX6SwJBsgwLmDOJHU7,iv:Fx1bpRKrz4wKztuEXMfa0KuRqLcOu9ZLT8OYdH+i58c=,tag:IoDhNZpGnGhqmDllgUVdUg==,type:str]" @secret(sops)
}

SECRET: "ENC[AES256_GCM,data:sNFoMGJGHOZZ0liZZ7HHtY/3BpJcQQ==,iv:JA22HeW9V00WEKINOGt/y5enl+94gOmZ6TvMlCYQZl4=,tag:iqGt5iVYqvMg+onsUXjAFQ==,type:str]" @secret(sops)

another: {
	this:  100
	value: "is not encrypted"
}

// DO NOT EDIT: auto-generated by cue-sops
sops: {
	kms:      null
	gcp_kms:  null
	azure_kv: null
	hc_vault: null
	age: [{
		recipient: "age1ethasxep4zkax64yfx35rn2t4yeul4254w764l9gtasvn2rwpv7s733dq7"
		enc: """
			-----BEGIN AGE ENCRYPTED FILE-----
			YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBQMk43bjFuUytFNTlIclNW
			Y1RnNWEwc2FGOUd4VW5NODNwdEVKOXlJd2k4ClNubDN0Qktuck5IVnN6ZjZBOTEz
			OFhNRUc3aUs1Y09DQTF6OTlWRU9ZQ00KLS0tIGdGTUlZWUJyVkZKZXdvMzZhV294
			c0E5bHVkSHc0MkhFUnhiODFlbzV5SE0KlNEhfwHl/VDZzfkpGb2/s7KbTFRA4U/K
			u5OM5P2YTvpSkmVbdVLLcX7eFHVyLZOukarFXEZ65rq9baMO0lJ3Vg==
			-----END AGE ENCRYPTED FILE-----

			"""
	}]
	lastmodified:    "2022-12-21T13:33:41Z"
	mac:             "ENC[AES256_GCM,data:heUT68PAirogTfcV+4pR8RNjx+d3cEE+Zn5e97xNy2wJvwZ4ecxnxItDj60E71aTK80UxCxkWkfjg2ZGKscPCMKoAXBkli6y/ab0e0+9uulvqjbd51m7mzGo/DMt65Ab7C6hq6S/VuI9JvvR7OVdgpvrliQzlCx2VENYNG6/r/0=,iv:gPvKgisLoTuOEIMNQgwY3zhPUEDkjJrRTyGWEEMr1ww=,tag:P8OlN/XDfWZqo6ZIchwbzw==,type:str]"
	pgp:             null
	encrypted_regex: "api_key|SECRET"
	version:         "3.7.3"
}
```

Let's also convert this to JSON:

```bash
cue val secrets.cue --out json > secrets.json
```

To decrypt the CUE file execute:

```
go run main.go decrypt secrets.cue
```

Our file now looks like this:

```
data: {
	api_key: "this is an api key" @secret(sops)
}

SECRET: "this is a secret value" @secret(sops)

another: {
	this:  100
	value: "is not encrypted"
}
```

Verify that we have interoperability with the sops CLI:

```shell
sops -d secrets.json
```

## Current Limitations

- We currently can't handle values in the file that are not concrete, for example constraints.

