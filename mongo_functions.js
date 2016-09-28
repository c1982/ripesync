
GetActiveIpCountByAsn
function(asnumber) {
    return db.hostactivity.aggregate([
    {$match: {asn:asnumber}},
    {$group: {_id: "$ip"}},
    {$group: { _id: 1, count: { $sum: 1}}}
    ]);
}

GetIPCountByPortNumber
function(asnumber, portnumber) {
    
return db.hostactivity.aggregate([
    {$match: {asn:asnumber, ports: {$elemMatch: {port: portnumber}}}},
    {$group: {_id: "$ip"}},
    {$group: { _id: 1, count: { $sum: 1}}}
    ]);
}

GetIPCountByTTLRange
function(asnumber, low, hight) {
    return db.hostactivity.aggregate([
    {$match: {asn:asnumber, ports: {$elemMatch: {ttl:{ $gte: low, $lt: hight }}}}},
    {$group: {_id: "$ip"}},
    {$group: { _id: 1, count: { $sum: 1}}}
    ]);
}