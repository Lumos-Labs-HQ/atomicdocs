import type { RequestHandler } from 'express';
import type { MiddlewareHandler } from 'hono';

interface ExpressApp {
  _router: {
    stack: Array<{
      route?: {
        path: string;
        methods: Record<string, boolean>;
        stack: Array<{ handle?: Function }>;
      };
    }>;
  };
  get(setting: string): unknown;
  set(setting: string, value: unknown): void;
}

interface HonoApp {
  routes: Array<{
    method: string;
    path: string;
    handler?: Function;
  }>;
}

/**
 * AtomicDocs middleware for Express.js and Hono frameworks.
 * Auto-detects the framework and returns the appropriate middleware.
 * 
 * @example
 * // Express.js usage
 * import express from 'express';
 * import atomicdocs from 'atomicdocs';
 * 
 * const app = express();
 * app.use(atomicdocs());
 * 
 * @example
 * // Hono usage
 * import { Hono } from 'hono';
 * import atomicdocs from 'atomicdocs';
 * 
 * const app = new Hono();
 * app.use('*', atomicdocs(app, 3000));
 */
declare function atomicdocs(): RequestHandler;
declare function atomicdocs(app: HonoApp, port: number): MiddlewareHandler;

declare namespace atomicdocs {
  /**
   * Manually register Express routes with AtomicDocs.
   * Call this after all routes are defined.
   * 
   * @param app - Express application instance
   * @param port - Port number the server is running on
   * 
   * @example
   * import express from 'express';
   * import atomicdocs from 'atomicdocs';
   * 
   * const app = express();
   * const PORT = 3000;
   * 
   * app.use(atomicdocs());
   * 
   * // Define routes...
   * 
   * app.listen(PORT, () => {
   *   atomicdocs.register(app, PORT);
   * });
   */
  function register(app: ExpressApp, port: number): void;
}

export = atomicdocs;
