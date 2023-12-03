If you don't have make install

# To install locally



You need also rust and its toolchain

```
sudo zypper install make


cargo install cargo-watch
cargo install --locked wasm-bindgen-cli
cargo install --locked trunk
```




```
cargo watch -w ./src  -s "trunk serve"
```


Links:
- https://www.youtube.com/watch?v=S-O9QkrlfYw
- good resources for components https://github.com/wiseaidev/yew-components-pure-css/tree/main/src/components
- flexbox grid https://codepen.io/mcraiganthony/pen/NxGxqm