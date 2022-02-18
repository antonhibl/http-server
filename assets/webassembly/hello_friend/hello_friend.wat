(;
Every WAT application must be a module. WAT doesn't have native string 
support, so working with strings requires you to work directly with memory 
as an array of character data. That memory data then must be converted into
a string in JavaScript code, because manipulating strings from within 
JavaScript is much simpler. When working with strings in WAT, you need to
declare an array of character data that is stored within WebAssembly linear
memory. I will also need to call an imported JavaScript function from 
WebAssembly to handle I/O operations. Unlike in a native application where
the operating system usually handles I/O. in a WebAssembly module, I/O must
be handled by the embedding environment, whether that environment is a web 
browser, an operating system, or runtime.
;)
(module 
    (import "env" "print_string" 
        (func $print_string( param i32)))
            (import "env" "buffer" (memory 1)) 
            (global $start_string
                (import "env" "start_string") i32) (global $string_len i32
                (i32.const 14))
                (data (global.get $start_string) "Hello, friend.")
                (func (export "hellofriend") 
                (call $print_string (global.get $string_len)) )
)