---

TODO:

1: support decryption

would be great to filter fields from the cue so that we could encrypt files using contraints

make re-encrpytion idempotent?

---

secrets plan

we encrypt secrets using sops

we add commands to the compont cli to encrypt or descrpyt secrets using SOPS library

secrets can then be stored in Git (can we make that bit easier?)

---

using sops like encryption:

we add to cue attributes, one for field level encryption another for file level

we using x509 certificates to encrypt the data; if we can use fulcio or SPIFFE to enable this that would be best.

whenever a component is built or rendered we use the available certs for enc/dec

we add an encrypt command that can be used to encrypt on disk (probably same for decrypt)

process:
- read cue file,
- look for @(secret) attribute
- get x509 certs
- encrypt fields data
- write file

maybe we can do this purely using the ast?

---

another plan would be to use SPIFFE to generate x509 certificates that sops could use encrpyt secrets.

users could authenticate using a SPIFFE ID and sign

in cluster we would decrypt using SVID, policies can then be used to determine which identities can decrypt secrets.

only problem is that sops does not support x509
