package research

const testBatch = `{
"title":"Meridian Archive","author":"","publisher":"","classification":"primary",
"entities":[{"id":"observatory","name":"Meridian Observatory","kind":"organization","summary":"Synthetic subject","time":{},"attributes":{},"version":1}],
"events":[{"id":"foundation","title":"Foundation","kind":"Formation","summary":"The synthetic observatory was established in 1800.","time":{"start":{"year":1800,"precision":"year","calendar":"proleptic_gregorian"}},"participants":["observatory"],"locationIds":[],"claimIds":["claim"],"importance":10,"phase":"Origins"}],
"claims":[{"id":"claim","text":"The synthetic observatory was established in 1800.","status":"unverified","eventId":"foundation","evidenceIds":["evidence"],"assessment":[]}],
"relationships":[],"conflicts":[],
"evidence":[{"id":"evidence","claimId":"claim","sourceId":"","excerpt":"The synthetic observatory was established in 1800.","location":"Foundation","relation":"supporting","mode":"quotation","verified":false}]
}`
