# GPIO setup

## Using Linux Device-Tree
For ring trigger it is necessary to configure a Key Stroke. This can be achived with configuring an gpio.

Copy the `gpio.dts` to `/boot/overlay-user/gpio.dts` and if you are using armbian execute
```
armbian-add-overlay /boot/overlay-user/gpio.dts
```

This will compile the `dts` file to a `dtbo` compiled one and will append the overlay to your `armbianEnv.txt