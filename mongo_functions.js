
GetActiveIpCountByAsn
function(asnumber) {
    return db.getCollection('hostactivity').distinct("ip",{asn:asnumber}).length
}

GetIPCountByPortNumber
function(asnumber, portnumber) {
    return db.getCollection('hostactivity').distinct("ip",{asn:asnumber, ports: {$elemMatch: {port: portnumber}}}).length
}

GetIPCountByTTLRange
function(asnumber, low, hight) {
 return db.getCollection('hostactivity').distinct("ip",{asn:asnumber, ports: {$elemMatch: {ttl:{ $gte: low, $lt: hight }}}}).length
}