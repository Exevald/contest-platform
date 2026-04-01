/* eslint-disable max-len,@typescript-eslint/consistent-type-definitions */

interface Window {
	cefclientSendQuery: (
		name: string,
		args: unknown[],
		resolve: (value: unknown) => void,
		reject: () => void,
	) => RequestId,

	cefclientSendCallbackResult: (
		externalRequestId: RequestId,
		externalRequestName: string,
		status: boolean,
		returnData: unknown,
	) => void,

	cefclientDispatch: (name: string, id: RequestId, args: unknown[]) => void,
}
