import AdminConfigHandler from 'server/handlers/api/admin/config';
import AdminReloadHandler from 'server/handlers/api/admin/reload';
import { MetaInfoHandler, MetaInfoLiveHandler } from 'server/handlers/api/meta/info';
import MetaTelemetryHandler from 'server/handlers/api/meta/telemetry';
import { SearchHandler, SearchLiveHandler } from 'server/handlers/api/search';
import SourceHandler from 'server/handlers/api/source';
import HealthHandler from 'server/handlers/meta/health';
import FallbackHandler from 'server/handlers/view/fallback';
import FrontendHandler from 'server/handlers/view/frontend';

export default [
  AdminConfigHandler,
  AdminReloadHandler,
  MetaInfoHandler,
  MetaInfoLiveHandler,
  MetaTelemetryHandler,
  SearchHandler,
  SearchLiveHandler,
  SourceHandler,
  HealthHandler,
  FrontendHandler,
  FallbackHandler,
];
