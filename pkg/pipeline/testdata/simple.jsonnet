local DefaultHttpStep() = {
  kind: 'HTTP',
  config: {
    default_timeout: 10000000000,
    default_method: 'GET',
  },
};

{
  steps: [
    DefaultHttpStep(),
    {
      kind: 'JSON',
      config: {
        template: '{"result":"{{ .foo }}"}',
      },
    },
  ],
}
