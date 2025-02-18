package web3ext

// EireneJs eirene related apis
const EireneJs = `
web3._extend({
	property: 'eirene',
	methods: [
		new web3._extend.Method({
			name: 'getSnapshot',
			call: 'eirene_getSnapshot',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getAuthor',
			call: 'eirene_getAuthor',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotProposer',
			call: 'eirene_getSnapshotProposer',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotProposerSequence',
			call: 'eirene_getSnapshotProposerSequence',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSnapshotAtHash',
			call: 'eirene_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getSigners',
			call: 'eirene_getSigners',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Method({
			name: 'getSignersAtHash',
			call: 'eirene_getSignersAtHash',
			params: 1
		}),
		new web3._extend.Method({
			name: 'getCurrentProposer',
			call: 'eirene_getCurrentProposer',
			params: 0
		}),
		new web3._extend.Method({
			name: 'getCurrentValidators',
			call: 'eirene_getCurrentValidators',
			params: 0
		}),
		new web3._extend.Method({
			name: 'getRootHash',
			call: 'eirene_getRootHash',
			params: 2,
		}),
		new web3._extend.Method({
			name: 'getVoteOnHash',
			call: 'eirene_getVoteOnHash',
			params: 4,
		}),
		new web3._extend.Method({
			name: 'sendRawTransactionConditional',
			call: 'eirene_sendRawTransactionConditional',
			params: 2,
			inputFormatter: [null]
		}),
	]
});
`
