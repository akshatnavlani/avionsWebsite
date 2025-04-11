import { createClient } from '@supabase/supabase-js'

const supabaseUrl = 'https://hldbpwyoykzmchgmpxlu.supabase.co'
const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImhsZGJwd3lveWt6bWNoZ21weGx1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDI2Mzg0MzIsImV4cCI6MjA1ODIxNDQzMn0.PnEJu5DzGGWUEOAU8kXJMaa8VnXxL-yzGtnFzsQ7bqY'

export const supabase = createClient(supabaseUrl, supabaseAnonKey) 