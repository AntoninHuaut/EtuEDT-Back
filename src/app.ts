import { connect } from '/sql/index.ts';

connect().then(() => import('/server.ts'));
