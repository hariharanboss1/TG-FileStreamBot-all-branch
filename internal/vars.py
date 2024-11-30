import sys
from os import environ
from dotenv import load_dotenv

load_dotenv()


class Var(object):
#auth
auth_channel = environ.get('AUTH_CHANNEL')
auth_grp = environ.get('AUTH_GROUP')
AUTH_CHANNEL = [int(auth_channel) for auth_channel in environ.get('AUTH_CHANNEL', '').split() if id_pattern.search(auth_channel)]
AUTH_GROUPS = [int(ch) for ch in auth_grp.split()] if auth_grp else None
