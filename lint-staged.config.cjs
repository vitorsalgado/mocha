'use strict'

module.exports = {
  '*.{md,json}': 'prettier --write --ignore-unknown',
  '*.go': ['make fmt', 'make vet'],
}
