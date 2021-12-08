```jsonc
{
    "description": "Gallia est omnis divisa in partes tres, quarum unam incolunt Belgae, aliam Aquitani, tertiam, qui ipsorum lingua Celtae, nostra Galli appellantur.",  # some words about that package
    "homepage": "https://github.com/egomobile/ego-cli#README", # url to a home with some information
    "maintainer": { # some information about the maintainer
        "name": "e.GO Mobile", # the (nick-)name
        "homepage": "https://github.com/egomobile", # if no homepage is available, maybe the user's github page
        "contacts": ["mailto:foo@example.com"] # (optional) list of contact addresses
    },
    "author": { # (optional) some information of the author
        "name": "", # the name of the author, (nick-)name or company name
        "homepage": "https://www.microsoft.com/" # if no homepage is available, maybe the user's github page
        "contacts": ["mailto:foo@example.com"] # (optional) list of contact addresses
    },
    "repositories": [ # (optional) list of one or more repositories
        {
            "type": "git",
            "url": "https://github.com/egomobile/ego-cli"
        }
    ],
    "sources": {
        # keys contain names of systems, separated by comma, if needed
        "debian,ubuntu": { # Ubuntu & Debian
            "x32,x64": {  # Intel compatible 32 and 64 bit systems
                "requirements": [ # list of requirements to install BEFORE package is installed
                    "apt://wget",
                    "apt://apt-transport-https"
                ],
                "source": "apt://code" # source
            }
        },
        "macos": { # MacOS
            "*": { # all systems / architectures
                "source": "brew://visual-studio-code"
            }
        },
        "windows": { # Microsoft Windows
            "x64": { # only 64 bit systems
                "source": "winget://Microsoft.VisualStudioCode"
            }
        }
    }
}
```