go run ../cmd/precompute \
  -graphql http://localhost:8080/otp/gtfs/v1 \
  -bucket weekday_08 \
  -out ../data/artifacts
