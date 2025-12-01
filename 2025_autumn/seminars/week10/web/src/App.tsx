import { useState } from 'react'
// IMPORTANT: Load polyfill first, then CSS, then generated code
import './goog-polyfill'
import './App.css'

// Now import the generated code
import { KeyValueServiceClient } from './generated/ServerServiceClientPb'

const ENVOY_URL = import.meta.env.VITE_ENVOY_URL || 'http://localhost:8080'

function App() {
  const [key, setKey] = useState('')
  const [value, setValue] = useState('')
  const [result, setResult] = useState<string>('')
  const [loading, setLoading] = useState(false)

  const client = new KeyValueServiceClient(ENVOY_URL, null, null)

  const handleSetValue = async () => {
    if (!key || !value) {
      setResult('Please enter both key and value')
      return
    }

    setLoading(true)
    setResult('')

    // Access protobuf classes from global scope at runtime
    // @ts-ignore
    const SetValueRequest = (window as any).proto.api.SetValueRequest;
    
    const request = new SetValueRequest()
    request.setKey(key)
    request.setValue(value)

    try {
      await new Promise<void>((resolve, reject) => {
        client.setValue(request, {}, (err, response) => {
          if (err) {
            reject(err)
          } else {
            resolve()
          }
        })
      })
      setResult(`✓ Successfully set key "${key}" to "${value}"`)
      setValue('')
    } catch (error: any) {
      setResult(`✗ Error: ${error.message || 'Failed to set value'}`)
    } finally {
      setLoading(false)
    }
  }

  const handleGetValue = async () => {
    if (!key) {
      setResult('Please enter a key')
      return
    }

    setLoading(true)
    setResult('')

    // Access protobuf classes from global scope at runtime
    // @ts-ignore
    const GetValueRequest = (window as any).proto.api.GetValueRequest;

    const request = new GetValueRequest()
    request.setKey(key)

    try {
      const response = await new Promise<string>((resolve, reject) => {
        client.getValue(request, {}, (err, response) => {
          if (err) {
            reject(err)
          } else {
            resolve(response.getValue())
          }
        })
      })
      setResult(`✓ Value for "${key}": "${response}"`)
      setValue(response)
    } catch (error: any) {
      if (error.code === 5) {
        setResult(`✗ Key "${key}" not found`)
      } else {
        setResult(`✗ Error: ${error.message || 'Failed to get value'}`)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app">
      <div className="container">
        <h1>🔑 Key-Value Store</h1>
        <p className="subtitle">gRPC-Web Client</p>

        <div className="form-group">
          <label htmlFor="key">Key:</label>
          <input
            id="key"
            type="text"
            value={key}
            onChange={(e) => setKey(e.target.value)}
            placeholder="Enter key"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="value">Value:</label>
          <input
            id="value"
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="Enter value (for Set) or leave empty (for Get)"
            disabled={loading}
          />
        </div>

        <div className="button-group">
          <button
            onClick={handleSetValue}
            disabled={loading || !key || !value}
            className="btn btn-primary"
          >
            {loading ? 'Setting...' : 'Set Value'}
          </button>
          <button
            onClick={handleGetValue}
            disabled={loading || !key}
            className="btn btn-secondary"
          >
            {loading ? 'Getting...' : 'Get Value'}
          </button>
        </div>

        {result && (
          <div className={`result ${result.startsWith('✓') ? 'success' : 'error'}`}>
            {result}
          </div>
        )}

        <div className="info">
          <p>Connected to: {ENVOY_URL}</p>
        </div>
      </div>
    </div>
  )
}

export default App

