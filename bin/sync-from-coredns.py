#!/usr/bin/python

# Copy the <middleware>/README.md to middleware/<mddleware>.md, add
# some Hugo meta data to let Hugo render it.

import datetime
import os
import re
import sys

def header(middleware, title, description, weight):
  h = """
+++
title = "%(title)s"
description = "%(description)s"
weight = "%(weight)s"
tags = [ "middleware", "%(middleware)s" ]
categories = [ "middleware" ]
date = "%(date)s"
+++
""" % {'title': title, 'description': description,
      'weight': weight, 'middleware': middleware,
      'date': str(datetime.datetime.utcnow().isoformat()) }
  return h

def parse(readme, middleware):
  file = open(readme)

  description, title, rest = '', '', ''
  
  state = 'TITLE'
  line = file.readline()
  while line:
    if state == 'TITLE':
      title = re.sub('# ?', '', line).rstrip()
      state = 'SKIP'
      line = file.readline()
      continue

    if state == 'SKIP':
      if line.strip() == '':
        line = file.readline()
        continue
      state = 'DESCRIPTION'

    if state == 'DESCRIPTION':
      if line.strip() == '':
        state = 'REST'
      else:
        description += line.rstrip() + " "

    if state == 'REST':
      rest += line

    line = file.readline()

  file.close()
  h = header(middleware, title, description.rstrip(), weight)
  return h + rest

weight=0
tags='middleware'
content='../content/middleware'

if not os.path.isdir(content):
  sys.exit("%s: Need to be run from the site's bin directory" % sys.argv[0])

if len(sys.argv) < 2:
  sys.exit("%s: Need containing directory of CoreDNS middleware" % sys.argv[0])


for middleware in os.listdir(sys.argv[1]):
  dir = os.path.join(sys.argv[1], middleware)
  readme = os.path.join(dir, "README.md")
  page = os.path.join(content, middleware) + ".md"
  if os.path.isdir(dir) and os.path.exists(readme):
    weight+=1
    print >> sys.stderr, readme + " --> " + page
    p = parse(readme, middleware)
    filep = open(page, 'w')
    filep.write(p)
    filep.close()
