const fastify = require("fastify")({ logger: true });

fastify.get("/hello", async (request, reply) => {
  return { hello: "world" };
});

fastify.listen(3000);
