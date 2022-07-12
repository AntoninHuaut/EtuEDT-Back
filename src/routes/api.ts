import { Context, Router } from 'oak';
import v2router from './v2.ts';

const router = new Router();

router.get('/', async (context: Context) => {
    context.response.body = {
        path: 'v2',
    };
});

v2router(router, '/v2');

export default router;
