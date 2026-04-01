/* eslint-disable @typescript-eslint/consistent-type-definitions, no-console, max-params, no-void */

type Handler = (...args: unknown[]) => unknown

type CefAdapter = {
	call: <T = unknown>(
		name: string,
		args?: unknown[],
		signal?: AbortSignal,
	) => Promise<T>,

	listen: (
		eventName: string,
		handler: Handler,
	) => void,

	unlisten: (
		eventName: string,
	) => void,
}

type RequestId = number

type CefConnector = {
	sendRequest: (
		name: string,
		args: unknown[],
		resolve: (value: unknown) => void,
		reject: (reason?: unknown) => void,
	) => RequestId,
	sendResponse: (
		externalRequestId: RequestId,
		externalRequestName: string,
		status: boolean,
		returnData: unknown,
	) => void,
}

const cefConnector: CefConnector = {
	sendRequest: (name, args, resolve, reject) => window.cefclientSendQuery(name, args, resolve, reject),
	sendResponse: (id, name, ok, data) => window.cefclientSendCallbackResult(id, name, ok, data),
}

const handleCefEvent = (
	name: string,
	id: RequestId,
	args: unknown[],
	handlers: Map<string, Handler>,
	connector: CefConnector,
): void => {
	const handler = handlers.get(name)
	if (!handler) {
		console.error(`CEF: no handler for event "${name}"`)
		connector.sendResponse(id, name, false, 'Handler not found')
		return
	}

	void (async () => {
		try {
			const result = await handler(...args)
			connector.sendResponse(id, name, true, result)
		}
		catch (err) {
			console.error(`CEF handler "${name}" failed:`, err)
			connector.sendResponse(id, name, false, err)
		}
	})()
}

const createCefAdapter = (connector: CefConnector = cefConnector): CefAdapter => {
	const handlers = new Map<string, Handler>()

	window.cefclientDispatch = (name, id, args) => {
		handleCefEvent(name, id, args, handlers, connector)
	}

	const call = <T = unknown>(
		name: string,
		args: unknown[] = [],
		signal?: AbortSignal,
	): Promise<T> => {
		if (signal?.aborted) {
			return Promise.reject('Request aborted')
		}

		return new Promise<T>((resolve, reject) => {
			connector.sendRequest(
				name,
				args,
				(value: unknown) => resolve(value as T),
				(reason?: unknown) => reject(reason),
			)

			const onCancel = () => {
				reject('Request aborted')
			}

			signal?.addEventListener('abort', onCancel, {once: true})
		})
	}

	const listen = (eventName: string, handler: Handler): void => {
		handlers.set(eventName, handler)
	}

	const unlisten = (eventName: string): void => {
		handlers.delete(eventName)
	}

	return {
		call,
		listen,
		unlisten,
	}
}

export {
	createCefAdapter,
	handleCefEvent,
	CefConnector,
	RequestId,
	Handler,
	CefAdapter,
}
