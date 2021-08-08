
# Load analysis 
# Generate cpu logs from top, using
# First compile with wInvalidate
#   pgrep 11_improved_animation
#   top -l 0 -s 1  -pid 8849 -stats cpu | awk 'NR%13==0; fflush(stdout)' > wInvalidate.txt

# recompile top opInvalidateOp
#   pgrep 11_improved_animation
#   top -l 0 -s 1  -pid 9144 -stats cpu | awk 'NR%13==0; fflush(stdout)' > opInvalidateOpAdd.txt

# Run Egg timer for 60 seconds

library(data.table)
library(ggplot2)

# Source
d1 = fread("opInvalidateOpAdd.txt")
d2 = fread("wInvalidate.txt")

# Rename
setnames(d1, "V1", "opInvalidateOpAdd")
setnames(d2, "V1", "wInvalidate")

# Combine and arrange
d = cbind(d1,d2)
d[, sec:=-2:(.N-3)]
md = melt(d, id.vars="sec")

# Visualize
p = ggplot(md, aes(x=sec, y=value, color=variable)) + 
    geom_line() +
    labs(x="Time (s)", 
         y="CPU load (%)", 
         title="Egg timer CPU load", 
         subtitle="Comparing invalidation methods", 
         caption="MacBook Air (2017) \n1,8 GHz Dual-Core Intel Core i5 \nAug 8th, 20201")

# Save
ggsave("../../11_invalidate_cpu_load.png", width=6, height=4, dpi="print")
