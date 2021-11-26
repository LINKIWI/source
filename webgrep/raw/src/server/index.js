import Express from 'express';
import { supercharge } from 'supercharged/server';
import path from 'path';
import yargs from 'yargs';
import * as Sentry from '@sentry/node';
import Context from 'server/context';
import handlers from 'server/handlers';

const { VERSION } = process.env;

const main = () => {
  const params = yargs
    .usage('$0 [options]', 'webgrep: web frontend for Livegrep code search')
    .option('config', {
      alias: 'c',
      default: 'config.yaml',
      description: 'Path to the configuration file on disk',
    })
    .option('verbosity', {
      alias: 'v',
      default: 'info',
      description: 'Log verbosity level',
      choices: ['error', 'warn', 'info', 'debug'],
    })
    .version(VERSION)
    .argv;

  const app = Express();
  const ctx = new Context({
    config: params.config,
    verbosity: params.verbosity,
  });

  ctx.log.info(
    'main: starting webgrep: version=%s config=%s verbosity=%s',
    VERSION,
    params.config,
    params.verbosity,
  );

  const sentryDSN = ctx.config.get('server.sentry_dsn');
  if (sentryDSN) {
    Sentry.init({ dsn: sentryDSN });
    Sentry.configureScope((scope) => {
      scope.setTag('version', VERSION);
    });
    ctx.log.debug('main: enabled Sentry reporting: dsn=%s', sentryDSN);
  }

  app.use(Sentry.Handlers.requestHandler());
  app.use('/assets', Express.static(path.join(__dirname, '../../dist/client')));
  supercharge(app, handlers, {
    createHandler: (HandlerClass) => (...args) => new HandlerClass(...args, ctx),
    ws: { perMessageDeflate: true },
  });
  app.use(Sentry.Handlers.errorHandler());

  const transport = ctx.config.get('server.listen.transport');
  const address = ctx.config.get('server.listen.address');
  switch (transport) {
    case 'tcp': {
      const [host, port] = address.split(':', 2);
      ctx.log.info('main: serving indefinitely over TCP: host=%s port=%d', host, port);
      app.listen(port, host);
      return 0;
    }
    case 'unix':
      ctx.log.info('main: serving indefinitely over Unix domain socket: socket=%s', address);
      app.listen(address);
      return 0;
    default:
      ctx.log.error('main: unsupported server listen transport: transport=%s', transport);
      return 1;
  }
};

const code = main();
if (code) {
  process.exit(code);
}
