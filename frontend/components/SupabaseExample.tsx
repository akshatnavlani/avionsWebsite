'use client'

import { useEffect, useState } from 'react'
import { supabase } from '@/lib/supabase'

export default function SupabaseExample() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function fetchData() {
      try {
        // Example query - replace 'your_table' with your actual table name
        const { data, error } = await supabase
          .from('your_table')
          .select('*')
          .limit(10)

        if (error) throw error

        setData(data || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error}</div>

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Supabase Data Example</h2>
      <div className="grid gap-4">
        {data.map((item, index) => (
          <div key={index} className="p-4 border rounded-lg">
            <pre className="whitespace-pre-wrap">
              {JSON.stringify(item, null, 2)}
            </pre>
          </div>
        ))}
      </div>
    </div>
  )
} 