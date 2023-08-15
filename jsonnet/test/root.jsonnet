local dz = import '../http.libsonnet';
local m = import '../matchers.libsonnet';

local req = dz.req;

dz.mock([
  dz.repeat(10),
  dz.scenario('scenario_test'),
  dz.request([
    req.scheme('http'),
    req.method('post'),
    req.header('hello', m.eq('world')),
    req.header('content-type', 'application/json'),
    req.header('test', m.present()),
    req.header('complex', m.all([m.eq('test'), m.present()])),
    req.query('filter', m.both(m.eq('test'), m.present())),
  ]),
  dz.postActions([dz.postAction('pa', { name: 'test' }), dz.postAction('no_args'), dz.postAction('basic', 'dev')]),
])
