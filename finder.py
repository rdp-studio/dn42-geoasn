import os
import constant
import re
import csv

pattern_route = re.compile(r'^\s*route:\s*(\S+)', re.IGNORECASE)
pattern_route6 = re.compile(r'^\s*route6:\s*(\S+)', re.IGNORECASE)
pattern_origin = re.compile(r'^\s*origin:\s*(\S+)', re.IGNORECASE)
pattern_as_name = re.compile(r'^\s*as-name:\s*(\S+)', re.IGNORECASE)

def find_as_name(asn):
    asn = str(asn).upper()
    if not asn.startswith("AS"):
        asn = "AS" + asn

    try:
        with open(os.path.join(constant.REGISTRY_PATH, "data/aut-num", asn), encoding="utf-8") as f:
            for line in f:
                if pattern_as_name.match(line):
                    return pattern_as_name.match(line).group(1)
            return None
    except:
        return None

def list_route_and_origin():
    routes = []
    f_list = os.listdir(constant.REGISTRY_PATH + "/data/route")

    for f in f_list:
        with open(constant.REGISTRY_PATH + "/data/route/" + f, encoding="utf-8") as f:
            route = None
            origin = None
            for line in f:
                if pattern_route.match(line):
                    route = pattern_route.match(line).group(1)
                elif pattern_origin.match(line):
                    origin = pattern_origin.match(line).group(1)
            if route and origin:
                routes.append((route, origin, find_as_name(origin)))
    return routes

def list_route6_and_origin():
    routes = []
    f_list = os.listdir(constant.REGISTRY_PATH + "/data/route6")

    for f in f_list:
        with open(constant.REGISTRY_PATH + "/data/route6/" + f, encoding="utf-8") as f:
            route = None
            origin = None
            for line in f:
                if pattern_route6.match(line):
                    route = pattern_route6.match(line).group(1)
                elif pattern_origin.match(line):
                    origin = pattern_origin.match(line).group(1)
            if route and origin:
                routes.append((route, origin, find_as_name(origin)))
    return routes

if __name__ == "__main__":
    route = list_route_and_origin()
    route6 = list_route6_and_origin()
    merge_route = route + route6

    csv_rows = []

    for r in merge_route:
        if r[2] == None:
            print(f"{r[0]} {r[1]} no name, skipping")
            continue
        
        csv_rows.append([r[0], r[1].replace("AS", ""), r[2]])
    
    with open(constant.SOURCE_OUTPUT, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerows(csv_rows)
