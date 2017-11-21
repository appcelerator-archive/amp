const swaggerCombine = require('swagger-combine');

swaggerCombine('root.swagger.yml')
    .then(res => console.log(JSON.stringify(res)))
    .catch(err => console.error(err));
