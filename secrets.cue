package main

data: {
	api_key: "ENC[AES256_GCM,data:IlDe1sx/OibYexp5spNF8AA=,iv:jI3mPI5An4bXaZn6089kN8CdUa7Xz77qVY+dBEbc8qY=,tag:ipFASthZNQQruA8PhOSHKQ==,type:str]" @secret(enc,kind=sops)
}

api_key: "ENC[AES256_GCM,data:ojefIqJslYiw9fQAhvqbp7A=,iv:tchb4QchaEqQdskQNfPDvKiquUtYcwnMtHVxqlb35uA=,tag:fFFkvZPC/kaw99VGC6spzw==,type:str]" @secret(enc,kind=sops)

another: {
	this: 100
}

// DO NOT EDIT: auto-generated by tooling
sops: {
	kms:      null
	gcp_kms:  null
	azure_kv: null
	hc_vault: null
	age: [{
		recipient: "age1h0nshge5hff7zhzd8rzr8jpl2xwsr3q9mtw8sy05z7gseswd35qq2ewegh"
		enc: """
			-----BEGIN AGE ENCRYPTED FILE-----
			YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBqdEwrR1NUa3dqV0lKMmhv
			RFVHNTBaeTNDelQvUTFLN1N0UE9TSlFwa3dNCmlVNjhvVWZxSFM4Ty9ZTThZNFFn
			a2VWWmFYb3BtdHRua0FZY0x4SmE1ZHMKLS0tIDZqd2NvdTZzOWtZWFpDUkY2Nmd3
			U3ZuSUF4SmVtU1EvUzZLWHRQcDJVOFkKwcQ9VSqYo8hVCjs9ArIpEBVH977glxR7
			2ca8o3ReHJljQpJbZGckvBzhJeCUMWa6Ue/9+qo4vZHdTLLcE5/srg==
			-----END AGE ENCRYPTED FILE-----

			"""
	}]
	lastmodified:    "2022-12-20T19:06:24Z"
	mac:             "ENC[AES256_GCM,data:KH9Q9dbS/YFU59pCqMh5BqiKlCwDojiu2QV31T0R+e77KV4/JBUQUmWpYDU9rlbPNfOBaSkK2iG7H0SLpyVJm+T/zp0m3l4rywGQAaa3JT9RkkH5/lGxtPdooQuN3ppCz27bGJn8xnErZ4NX1tVJrfIGacm7LADB3pnqKyAO9PI=,iv:RnUnrR465e8HXCVQAUPmYMyfNa7AQlp1za2+5+RBMt8=,tag:HpAI9zVyYY2Rb0KqwwIbxQ==,type:str]"
	pgp:             null
	encrypted_regex: "api_key|api_key"
	version:         "3.7.3"
}
