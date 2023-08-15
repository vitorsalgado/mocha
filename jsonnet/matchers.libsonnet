local matcher(name, parameters) = [name, parameters];
local no_args_matcher(name) = [name];

{
  all(v):: matcher('all', v),
  any(v):: matcher('any', v),
  anything():: no_args_matcher('anything'),
  both(a, b):: matcher('both', [a, b]),
  eq(v):: matcher('eq', v),
  present():: no_args_matcher('present'),
}
