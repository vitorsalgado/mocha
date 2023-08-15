{
  merge(objs)::
    local merge(arr, i, running) =
      if i >= std.length(arr) then running
      else merge(arr, i + 1, std.mergePatch(running, arr[i])) tailstrict;
    merge(objs, 0, {}),
}
