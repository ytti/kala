# Kala
Library and CLI UI for managing passwords. I have no clue about crypto and I
only wrote this to learn basics of golang. So password manager by security and
programming language incompetent developer, you're gonna love this.

## Security
Every entry has 'secret entry' which right now includes only 'passphrase' but
could have something like HOTP master key etc in future. This 'secret entry' is
encrypted with salsa20 and key is derived with scrypt from master passphrase
salted with the entry name.

Entry is put into array of entries, entries is encrypted with salsa20  and key
is derived with scrypt from master password with compile time static salt.

Runtime only secrets which you're interested in are decrypted, secrets for
entries you didn't ask for are not touched. So perhaps in some corner-cases
when memory is compromised but master password is not, you avoid compromising
all of the passwords. I don't know if this makes sense or not.

I tried not to copy key or masterpassword in memory, rather just pass pointer
(slice when copied does not copy content, unlike array, so password slice is
usually copied straight off). Not sure if that makes any sense either.

## Install
go get https://github.com/ytti/kala/

## Use
```
[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala
Creating new /home/ytti/.config/kala/kala.json, new master passphrase needed
          Passphrase [] : 
  Passphrase (again) [] : 
Empty /home/ytti/.config/kala/kala.json created
[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala -a "google apps"
Passphrase for /home/ytti/.config/kala/kala.json: 
     Name [google apps] : 
     Host [google apps] : https://mail.ytti.fi
            Username [] : saku@ytti.fi
         Information [] : 
          Passphrase [] : 
  Passphrase (again) [] : 
Add Entry (YES/no): yes
Entry written to file
[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala -a "twitter"    
Passphrase for /home/ytti/.config/kala/kala.json: 
         Name [twitter] : 
         Host [twitter] : https://twitter.com
            Username [] : ytti
         Information [] : nada
          Passphrase [] : 
  Passphrase (again) [] : 
Add Entry (YES/no): yes
Entry written to file
[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala --export
Passphrase for /home/ytti/.config/kala/kala.json: 
[
  {
    "Name": "google apps",
    "Host": "https://mail.ytti.fi",
    "Username": "saku@ytti.fi",
    "Information": "",
    "Secret": {
      "Passphrase": "poop"
    }
  },
  {
    "Name": "twitter",
    "Host": "https://twitter.com",
    "Username": "ytti",
    "Information": "nada",
    "Secret": {
      "Passphrase": "meh"
    }
  }
]
[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala twit    
Passphrase for /home/ytti/.config/kala/kala.json: 
twitter
        Host : https://twitter.com
    Username : ytti
 Information : nada
  Passphrase : meh
[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala ytti
Passphrase for /home/ytti/.config/kala/kala.json: 
google apps
        Host : https://mail.ytti.fi
    Username : saku@ytti.fi
 Information : 
  Passphrase : poop

twitter
        Host : https://twitter.com
    Username : ytti
 Information : nada
  Passphrase : meh

[1 ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala -e twitter
Passphrase for /home/ytti/.config/kala/kala.json: 

twitter
        Host : https://twitter.com
    Username : ytti
 Information : nada
Edit this entry (yes/NO): yes
         Name [twitter] : 
Host [https://twitter.com] : http://twitter.com
        Username [ytti] : 
     Information [nada] : 
          Passphrase [] : 
  Passphrase (again) [] : 
Entry changed

1 entries changed, commit changes (yes/NO): yes
Changes written to file

[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% ./kala -d twitter
Passphrase for /home/ytti/.config/kala/kala.json: 

twitter
        Host : http://twitter.com
    Username : ytti
 Information : nada
Delete this entry (yes/NO): yes
Marking for deletion

1 entries marked for deletion, commit changes (yes/NO): yes
Deletions written to file

[ytti@ytti.fi ~/go/src/github.com/ytti/kala]% 
```
