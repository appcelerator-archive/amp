import babel from 'rollup-plugin-babel'
import eslint from 'rollup-plugin-eslint'
import svelte from 'rollup-plugin-svelte'
import commonjs from 'rollup-plugin-commonjs'
import nodeResolve from 'rollup-plugin-node-resolve'

export default {
  entry: 'src/main.js',
  dest: 'dist/main.js',
  format: 'es',
  plugins: [
    eslint({
      include: [
        './src/**/*.js',
      ],
    }),
    svelte({
      include: 'src/components/**.html',
    }),
    nodeResolve(),
    commonjs(),
    babel(),
  ],
}
