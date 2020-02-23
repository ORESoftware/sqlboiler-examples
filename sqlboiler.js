console.log(
  JSON.stringify({
    pkgname: 'model',
    output: 'v1/model',
    psql: {
      dbname: process.env.db_name || 'postgres',
      host: process.env.db_host || 'localhost',
      user: process.env.db_user || 'postgres',
      pass: process.env.db_pwd || 'postgres',
      blacklist: []
    },
    types: [
      {
        match: {
          type: 'types.Decimal'
        },
        replace: {
          type: 'apd.Decimal'
        },
        imports: {
          third_party: ['github.com/cockroachdb/apd']
        }
      },
      {
        match: {
          db_type: 'uuid',
          nullable: false
        },
        replace: {
          type: 'uuid.UUID'
        },
        imports: {
          third_party: ['github.com/google/uuid']
        }
      }
    ]
  }, null, 2)
);
