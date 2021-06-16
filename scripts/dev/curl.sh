
####################
# Trade

# trade templates
curl -i -H "Accept: application/json" -H "X-AUTH-TOKEN: 1"  http://localhost:8000/v1/trades/tradeTemplates

# stage upload / download

curl -i -H "Accept: application/json" -H "X-AUTH-TOKEN: 1"  http://localhost:8000/v1/trades/trade-docs/123

curl -F "formfile=@/home/robert/test.pdf" -F 'stageIdx=1' -F 'tid=799249' -F note="my description" -F expiresAt="2006-01-02T15:04:05Z07:00" localhost:8000/v1/trades/trade-docs
