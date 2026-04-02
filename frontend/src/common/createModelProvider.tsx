import {
	action,
	atom,
	computed,
} from '@reatom/core'
import type {
	Atom, AtomLike, Fn,
} from '@reatom/core'
import {
	useContext, useRef, createContext,
} from 'react'
import type {
	PropsWithChildren,
	ComponentType,
} from 'react'

/**
 * Обёртка над atom для удобного создания атомов в модели
 */
function defAtom<T>(initialState: T, name?: string): Atom<T, [T | ((prev: T) => T)]> {
	return atom(initialState, name)
}

/**
 * Обёртка над computed для удобного создания вычисляемых атомов
 */
function defComputed<T>(fn: () => T, name?: string) {
	return computed(fn, name)
}

/**
 * Обёртка над action для удобного создания экшенов
 */
function defAction<T extends Fn>(fn: T, name?: string) {
	return action(fn, name)
}

/**
 * Результат создания провайдера модели
 */
type ModelProviderResult<ARGS, MODEL extends Record<string, AtomLike>> = {
	Provider: ComponentType<PropsWithChildren<ARGS>>,
	useModel: () => MODEL,
}

/**
 * Создаёт провайдер модели с контекстом и хуком
 */
function createModelProvider<ARGS, MODEL extends Record<string, AtomLike>>(
	createModel: (props: ARGS) => MODEL,
): ModelProviderResult<ARGS, MODEL> {
	const ModelContext = createContext<MODEL | null>(null)

	const Provider = (args: PropsWithChildren<ARGS>) => {
		const {children, ...props} = args
		const modelRef = useRef<MODEL | null>(null)

		if (!modelRef.current) {
			modelRef.current = createModel(props as ARGS)
		}

		return (
			<ModelContext.Provider value={modelRef.current as MODEL | null}>
				{children}
			</ModelContext.Provider>
		)
	}

	const useModel = () => {
		const ctx = useContext(ModelContext)
		if (!ctx) {
			throw new Error('useModel must be used within a Provider')
		}
		return ctx
	}

	return {
		Provider,
		useModel,
	}
}

export {
	defAtom,
	defComputed,
	defAction,
	createModelProvider,
}
