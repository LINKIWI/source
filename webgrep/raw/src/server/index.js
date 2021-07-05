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

  ctx.log.info('main: starting webgrep');

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

  const [host, port] = ctx.config.get('server.listen_addr').split(':', 2);
  ctx.log.info('main: serving indefinitely: host=%s port=%d', host, port);
  app.listen(port, host);
};

main();
