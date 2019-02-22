# Helper script to fix meta.json; add resultPage and count.
# First cmd param: what was passed to index-generator as inxpath. ex.: ../website/generated/inx

import os
import sys
import io
import json

def count_torrents_in_index(inxpath, blocksize):
    files = os.listdir(os.path.dirname(inxpath))

    def filterfiles(item):
        return "inx" in item

    inxfiles = filter(filterfiles,files)
    return sum(1 for _ in inxfiles)*blocksize
        

meta = json.load(io.open(sys.argv[1]+".meta.json","r"))
meta["resultPage"] = "resultpage"
meta["entries"] = count_torrents_in_index(sys.argv[1],1000)
meta["inxUrlBase"] = "website/generated/inx"
meta["invUrlBase"] = "website/generated/inv"

json.dump(meta,io.open(sys.argv[1]+".meta.json", "w"))
